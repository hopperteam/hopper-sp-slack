# hopper-sp-slack
## Design
### Workflow
#### On Startup  
1. Register Service Provider  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/RegisterServiceProvider.svg "Register Service Provider")  
2. Get Workspace Data
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/GetWorkspaceData.svg "Get Workspace Data")   
#### Listening To Events  
1. User Subscription (home tab unsubscibed view)  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/UserSubscription.svg "User Subscription")  
2. User Log Off (home tab subscibed view)  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/UserLogOff.svg "User Log Off")   
3. Message Received  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/MessageReceived.svg "Message Received")  
4. Change User Data  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/ChangeUserData.svg "Change User Data")  
5. New User Joined Workspace  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/NewUserJoinedWorkspace.svg "New User Joined Workspace")  
6. Channel Name Changed  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/ChannelNameChanged.svg "Channel Name Changed")  
7. Channel Created  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/ChannelCreated.svg "Channel Created")  
8. Channel Deleted  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/ChannelDeleted.svg "Channel Deleted")  
9. User Joined Channel  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/UserJoinedChannel.svg "User Joined Channel")  
10. User Left Channel  
![alt text](https://raw.githubusercontent.com/hopperteam/hopper-sp-slack/master/doc/UserLeftChannel.svg "User Left Channel")  

### API
#### Slack Home Tab
https://api.slack.com/surfaces/tabs/using
#### Slack Event API
Get limitation warning: https://api.slack.com/events/app_rate_limited  
Handle messages tutorial: https://slack.dev/node-slack-sdk/events-api  
User data update event: https://api.slack.com/events/user_change  
User joined workspace event: https://api.slack.com/events/team_join  
Channel name changed event: https://api.slack.com/events/channel_rename  
Channel created event: https://api.slack.com/events/channel_created  
Channel deleted event: https://api.slack.com/events/channel_deleted  
User joined channel: https://api.slack.com/events/member_joined_channel  
User left channel: https://api.slack.com/events/member_left_channel  
#### Slack Web API  
Get list of users: https://api.slack.com/methods/users.list  
Get list of channels: https://api.slack.com/methods/conversations.list  
Get members of channels: https://api.slack.com/methods/conversations.members  
Redirects and linking: https://api.slack.com/reference/deep-linking    
#### Hopper API  
https://developer.hoppercloud.net/

### Types

## Deployment
