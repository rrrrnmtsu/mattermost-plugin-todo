// Local integration patch: enforce source-post access before Todo mutations.
package main

import "github.com/mattermost/mattermost-server/v5/model"

func (p *Plugin) canLinkPost(userID, postID string) bool {
	post, appErr := p.API.GetPost(postID)
	if appErr != nil || post == nil || post.DeleteAt != 0 {
		return false
	}
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil || channel == nil || channel.DeleteAt != 0 {
		return false
	}
	member, appErr := p.API.GetChannelMember(post.ChannelId, userID)
	if appErr != nil || member == nil {
		return false
	}
	return p.API.HasPermissionToChannel(userID, post.ChannelId, model.PERMISSION_READ_CHANNEL) &&
		p.API.HasPermissionToChannel(userID, post.ChannelId, model.PERMISSION_CREATE_POST)
}
