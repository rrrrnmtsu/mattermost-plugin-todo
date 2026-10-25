package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func TestPostLinkAuthorization(t *testing.T) {
	for stage, name := range []string{"missing_post", "deleted_post", "missing_channel", "archived_channel", "not_member", "no_read_permission", "no_post_permission", "allowed_member"} {
		t.Run(name, func(t *testing.T) {
			api := &plugintest.API{}
			userID, postID, channelID := model.NewId(), model.NewId(), model.NewId()
			post := &model.Post{Id: postID, ChannelId: channelID}
			channel := &model.Channel{Id: channelID}
			missing := model.NewAppError("authorization_test", "missing", nil, "", http.StatusNotFound)
			if stage == 0 {
				api.On("GetPost", postID).Return((*model.Post)(nil), missing).Once()
			} else {
				if stage == 1 {
					post.DeleteAt = 1
				}
				api.On("GetPost", postID).Return(post, nil).Once()
				if stage >= 2 {
					if stage == 2 {
						api.On("GetChannel", channelID).Return((*model.Channel)(nil), missing).Once()
					} else {
						if stage == 3 {
							channel.DeleteAt = 1
						}
						api.On("GetChannel", channelID).Return(channel, nil).Once()
						if stage >= 4 {
							if stage == 4 {
								api.On("GetChannelMember", channelID, userID).Return((*model.ChannelMember)(nil), missing).Once()
							} else {
								api.On("GetChannelMember", channelID, userID).Return(&model.ChannelMember{ChannelId: channelID, UserId: userID}, nil).Once()
								api.On("HasPermissionToChannel", userID, channelID, model.PERMISSION_READ_CHANNEL).Return(stage >= 6).Once()
								if stage >= 6 {
									api.On("HasPermissionToChannel", userID, channelID, model.PERMISSION_CREATE_POST).Return(stage == 7).Once()
								}
							}
						}
					}
				}
			}
			p := &Plugin{}
			p.SetAPI(api)
			assert.Equal(t, stage == 7, p.canLinkPost(userID, postID))
			api.AssertExpectations(t)
		})
	}
}

func TestAddRejectsInaccessiblePostBeforeTaskMutation(t *testing.T) {
	api := &plugintest.API{}
	postID := model.NewId()
	api.On("GetPost", postID).Return((*model.Post)(nil), model.NewAppError("authorization_test", "missing", nil, "", 404)).Once()
	p := &Plugin{}
	p.SetAPI(api)
	p.initializeAPI()
	request := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString(`{"message":"task","post_id":"`+postID+`"}`))
	request.Header.Set("Mattermost-User-ID", model.NewId())
	response := httptest.NewRecorder()
	// listManager is intentionally nil: validation must happen before any state write.
	p.ServeHTTP(nil, response, request)
	assert.Equal(t, http.StatusNotFound, response.Code)
	api.AssertExpectations(t)
}

func TestThreadReplyRechecksCurrentMembership(t *testing.T) {
	api := &plugintest.API{}
	userID, postID, channelID := model.NewId(), model.NewId(), model.NewId()
	api.On("GetPost", postID).Return(&model.Post{Id: postID, ChannelId: channelID}, nil).Once()
	api.On("GetChannel", channelID).Return(&model.Channel{Id: channelID}, nil).Once()
	api.On("GetChannelMember", channelID, userID).Return((*model.ChannelMember)(nil), model.NewAppError("authorization_test", "missing", nil, "", 404)).Once()
	p := &Plugin{}
	p.SetAPI(api)
	p.postReplyIfNeeded(userID, postID, "completed", "task")
	api.AssertNotCalled(t, "CreatePost")
	api.AssertExpectations(t)
}
