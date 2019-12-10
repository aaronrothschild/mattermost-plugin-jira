---
description: Get started by installing the Jira plugin from the Marketplace
---

# Installation \(Steps 1-2\)

## Step 1: Installing the Plugin on Mattermost 

#### Requirements

* For Jira 2.1 Mattermost server v5.14+ is required \(Certain plugin APIs became available\)
* For Jira 2.0, Mattermost Server v5.12+ is required

### Install via Plugin Marketplace \(Recommended\)

1. Go to **Main Menu &gt; Plugin** Marketplace in Mattermost
2. Search for "Jira" or manually find the plugin from the list and click **Install**
3. After the plugin has downloaded and been installed, click the **Configure** button
4. Go to **Plugins Marketplace &gt; Jira**
   1. Click the **Configure** button
   2. Generate a **Secret** for `Webhook Secret` and `Stats API Secret`
5. Go to the top of the screen and set **Enable Plugin** to `True`and then click **Save** to enable the Jira plugin.

### \(Alternative\) Install via Manual Upload

If your server doesn't have access to the internet, you can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-jira/releases) and upload it to your server via **System Console &gt; Plugin Management.** The binary releases on the page above, are the same as used by the Marketplace.

## Step 2: Install App on Jira Server

For Jira Server or Data Center instances, type in  `/jira install server <your-jira-url>` to any Mattermost channel as a Mattermost System Admin, and follow the steps posted to the channel. For Jira Cloud, type in `/jira install cloud <your-jira-url>`

You will be presented with the following instructions to get an app installed on the Jira server:

1. In Jira, Navigate to [**Settings &gt; Apps &gt; Manage Apps**](http://mmtest.atlassian.net/plugins/servlet/upm?source=side_nav_manage_addons).
   * For older versions of Jira, navigate to **Administration &gt; Applications &gt; Add-ons &gt; Manage add-ons**.
2. Click **Settings** at bottom of page, enable development mode, and apply this change.
   * Enabling development mode in Jira allows you to install apps that are not from the Atlassian Marketplace.
3. Click **Upload app**.
4. In the **From this URL field**, enter The URL returned in your Mattermost console - it is custom for your installation. 
5. Wait for the app to install. Once completed, you should see an "Installed and ready to go!" message in Jira.

Return to Mattermost and Use the `/jira connect` command to connect a Mattermost account with a Jira account.  This needs to be done by each Mattermost/Jira user. 

Once an account is connected - you can now use the "More Actions" \(...\) option of any message in the channel \(available when you hover over a message\) to create and attach issues in Jira as well as slash commands. 

