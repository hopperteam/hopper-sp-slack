# hopper-sp-slack
## Design
### Workflow
#### On Startup
1. Register Service Provider
2. Get Workspace Data
#### Listen for Events
1. User Subscription
2. Message Received  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/MessageReceived.svg "Message Received")
3. Change User Data
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/ChangeUserData.svg "Change User Data")
4. New User joined
5. User Left
6. Channel Name Changed  
![alt text](https://github.com/hopperteam/hopper-sp-slack/blob/master/ChannelNameChanged.svg "Channel Name Changed")
7. Channel Added
8. Channel Deleted

### API
#### Slack Event API
Get limitation warning: https://api.slack.com/events/app_rate_limited
#### Slack Web API  
Get list of users: https://api.slack.com/methods/users.list
Get list of channels: https://api.slack.com/methods/conversations.list
#### Hopper API  
https://developer.hoppercloud.net/

## Deployment
