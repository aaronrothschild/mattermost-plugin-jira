---
description: Setup notifications to keep a channel up to date on certain events in Jira
---

# Notification Subscriptions

### About Notifications

Notifications are driven by events that happen in Jira.  Whenever a user updates an issue in Jira, that event is sent over to Mattermost via . webhook - to the Jira plugin.  The jira plugin looks through its `Notification Subscription Rules` looking for events that meet crtieria for notifying a channel in Mattermost.  If no subscription rules are matched - the webhook event is dropped.  

### Who can setup Notification Subscriptions on a Mattermost channel?

Depending on your server's configuration - you may need to be a system administrator in order to setup notification subscriptions in a channel.  This is configured on the [plugin configuration](../setup/configuration.md#step-2-configure-webhooks-in-jira) screen.

### How can I setup a notification subscription for a channel?

Begin by typing `/jira subscribe` within the channel you want to set the notifications for.  You will then be presented with a modal screen to configure the subscription. 

### What are the events that can trigger a notification?

The basic events that can be tracked are:

* Updating an issue
* Deleting an Issue
* Commenting on an issue
* ....

### Which issue types can trigger a notification?

Any issue type in Jira can optionally trigger a notification.  You need to specify which issue types should trigger a subscription notification.

### How can I prevent some events from being posted?

Sometimes there are events that happen very often or are of little consequence and don't require notifying the channel of a change in Jira.  Filters can be setup to prevent notifications that meet certain criteria - such as having a particular value on a field. 

### Can I post a JQL Query to create a subscription?

Currently, we don't support this functionality because of the difficulty in parsing the JQL queries consistently has proven difficult.  You can still setup a webhook in Jira and point it directly at a Mattermost channel.

### Why is the JQL Query "approximate"?

There are certain situations that will require the JQL approximation to be 'massaged' before pasting it into Jira's Query box directly.  If your query contains custom fields with multiple words, they will need to be placed within quotes for example.  





 

