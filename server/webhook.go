// See License for license information.
// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.

package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
)

type Webhook interface {
	Events() StringSet
	PostToChannel(p *Plugin, channelId, fromUserId string) (*model.Post, int, error)
	PostNotifications(p *Plugin) ([]*model.Post, int, error)
}

type webhookField struct {
	name string
	id   string
	from string
	to   string
}

type webhook struct {
	*JiraWebhook
	eventTypes    StringSet
	headline      string
	text          string
	fields        []*model.SlackAttachmentField
	notifications []webhookNotification
	fieldInfo     webhookField
}

type webhookNotification struct {
	jiraUsername  string
	jiraAccountID string
	message       string
	postType      string
	commentSelf   string
}

func (wh *webhook) Events() StringSet {
	return wh.eventTypes
}

func (wh webhook) PostToChannel(p *Plugin, channelId, fromUserId string) (*model.Post, int, error) {
	if wh.headline == "" {
		return nil, http.StatusBadRequest, errors.Errorf("unsupported webhook")
	}

	post := &model.Post{
		ChannelId: channelId,
		UserId:    fromUserId,
		// Props: map[string]interface{}{
		// 	"from_webhook":  "true",
		// 	"use_user_icon": "true",
		// },
	}
	if wh.text != "" || len(wh.fields) != 0 {
		// Get instance for replacing accountids in text. If no instance is available, just skip it.
		ji, err := p.currentInstanceStore.LoadCurrentJIRAInstance()
		if err == nil {
			wh.text = replaceJiraAccountIds(ji, wh.text)
		}

		model.ParseSlackAttachment(post, []*model.SlackAttachment{
			{
				// TODO is this supposed to be themed?
				Color:    "#95b7d0",
				Fallback: wh.headline,
				Pretext:  wh.headline,
				Text:     wh.text,
				Fields:   wh.fields,
			},
		})
	} else {
		post.Message = wh.headline
	}

	_, appErr := p.API.CreatePost(post)
	if appErr != nil {
		return nil, appErr.StatusCode, appErr
	}

	return post, http.StatusOK, nil
}

func (wh *webhook) PostNotifications(p *Plugin) ([]*model.Post, int, error) {
	if len(wh.notifications) == 0 {
		return nil, http.StatusOK, nil
	}

	// We will only send webhook events if we have a connected instance.
	ji, err := p.currentInstanceStore.LoadCurrentJIRAInstance()
	if err != nil {
		// This isn't an internal server error. There's just no instance installed.
		return nil, http.StatusOK, nil
	}

	posts := []*model.Post{}
	for _, notification := range wh.notifications {
		var mattermostUserId string
		var err error

		// prefer accountId to username when looking up UserIds
		if notification.jiraAccountID != "" {
			mattermostUserId, err = p.userStore.LoadMattermostUserId(ji, notification.jiraAccountID)
		} else {
			mattermostUserId, err = p.userStore.LoadMattermostUserId(ji, notification.jiraUsername)
		}
		if err != nil {
			continue
		}

		// Check if the user has permissions.
		jiraUser, err2 := p.userStore.LoadJIRAUser(ji, mattermostUserId)
		if err2 != nil {
			// Not connected to Jira, so can't check permissions
			continue
		}
		client, err2 := ji.GetClient(jiraUser)
		if err2 != nil {
			p.errorf("PostNotifications: error while getting jiraClient, err: %v", err2)
			continue
		}
		// If this is a comment-related webhook, we need to check if they have permissions to read that.
		// Otherwise, check if they can view the issue.

		isCommentEvent := wh.Events().Intersection(commentEvents).Len() > 0
		if isCommentEvent {
			err = client.RESTGet(notification.commentSelf, nil, &struct{}{})
		} else {
			_, err = client.GetIssue(wh.Issue.ID, nil)
		}
		if err != nil {
			p.errorf("PostNotifications: failed to get self: %v", err)
			continue
		}

		notification.message = replaceJiraAccountIds(ji, notification.message)

		post, err := ji.GetPlugin().CreateBotDMPost(ji, mattermostUserId, notification.message, notification.postType)
		if err != nil {
			p.errorf("PostNotifications: failed to create notification post, err: %v", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}

func newWebhook(jwh *JiraWebhook, eventType string, format string, args ...interface{}) *webhook {
	return &webhook{
		JiraWebhook: jwh,
		eventTypes:  NewStringSet(eventType),
		headline:    jwh.mdUser() + " " + fmt.Sprintf(format, args...) + " " + jwh.mdKeySummaryLink(),
	}
}

func (p *Plugin) GetWebhookURL(teamId, channelId string) (string, error) {
	cf := p.getConfig()

	team, appErr := p.API.GetTeam(teamId)
	if appErr != nil {
		return "", appErr
	}

	channel, appErr := p.API.GetChannel(channelId)
	if appErr != nil {
		return "", appErr
	}

	v := url.Values{}
	secret, _ := url.QueryUnescape(cf.Secret)
	v.Add("secret", secret)
	v.Add("team", team.Name)
	v.Add("channel", channel.Name)
	return p.GetPluginURL() + "/" + routeIncomingWebhook + "?" + v.Encode(), nil
}
