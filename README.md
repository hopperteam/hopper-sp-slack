# hopper-sp-slack
## Design
### Workflow
#### On Startup
1. Register Service Provider
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/RegisterServiceProvider.svg "Register Service Provider")  
2. Get Workspace Data
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/GetWorkspaceData.svg "Get Workspace Data")   
#### Listening To Events  
1. User Subscription (home tab unsubscibed view)  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/UserSubscription.svg "User Subscription")  
2. User Log Off (home tab subscibed view)  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/UserLogOff.svg "User Log Off")   
3. Message Received  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/MessageReceived.svg "Message Received")
4. Change User Data
5. New User joined workspace
6. User Left workspace (delete corresponding subscription)
7. Channel Name Changed  
8. Channel Added
9. Channel Deleted
10. User joined channel
11. User left channel

### API
#### Slack Home Tab
https://api.slack.com/surfaces/tabs/using
#### Slack Event API
Get limitation warning: https://api.slack.com/events/app_rate_limited  
Handle messages tutorial: https://slack.dev/node-slack-sdk/events-api
#### Slack Web API  
Get list of users: https://api.slack.com/methods/users.list  
Get list of channels: https://api.slack.com/methods/conversations.list
Get members of channels: https://api.slack.com/methods/conversations.members
Redirects and linking: https://api.slack.com/reference/deep-linking  
#### Hopper API  
https://developer.hoppercloud.net/

### Types

## Deployment
