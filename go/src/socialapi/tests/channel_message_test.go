package main

import (
	"net/http"
	"socialapi/models"
	"socialapi/request"
	"socialapi/rest"
	"socialapi/workers/common/tests"
	"testing"

	"github.com/koding/runner"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChannelMessage(t *testing.T) {
	tests.WithRunner(t, func(r *runner.Runner) {
		Convey("While testing channel messages given a channel", t, func() {

			account, groupChannel, groupName := models.CreateRandomGroupDataWithChecks()

			nonOwnerAccount := models.CreateAccountInBothDbsWithCheck()

			nonOwnerSes, err := models.FetchOrCreateSession(nonOwnerAccount.Nick, groupName)
			So(err, ShouldBeNil)

			ses, err := models.FetchOrCreateSession(account.Nick, groupName)
			So(err, ShouldBeNil)

			Convey("message should be able added to the group channel", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)
				So(post.Id, ShouldNotEqual, 0)
				So(post.Body, ShouldNotEqual, "")
				Convey("message can be edited by owner", func() {

					initialPostBody := post.Body
					post.Body = "edited message"

					editedPost, err := rest.UpdatePost(post, ses.ClientId)
					So(err, ShouldBeNil)
					So(editedPost, ShouldNotBeNil)
					// body should not be same
					So(initialPostBody, ShouldNotEqual, editedPost.Body)
				})

				// for now social worker handles this issue
				Convey("message can be edited by an admin", nil)
				Convey("message can not be edited by non-owner", nil)

			})

			Convey("topic messages initialChannelId must be set as owner group channel id", func() {
				ses, err := models.FetchOrCreateSession(account.Nick, groupName)
				So(err, ShouldBeNil)
				So(ses, ShouldNotBeNil)

				topicChannel, err := rest.CreateChannelByGroupNameAndType(account.Id, "koding", models.Channel_TYPE_TOPIC, ses.ClientId)
				So(err, ShouldBeNil)
				So(topicChannel, ShouldNotBeNil)

				post, err := rest.CreatePost(topicChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)
				So(post.InitialChannelId, ShouldNotEqual, topicChannel.Id)
				publicChannel, err := rest.GetChannelWithToken(post.InitialChannelId, ses.ClientId)
				So(err, ShouldBeNil)
				So(publicChannel.TypeConstant, ShouldEqual, models.Channel_TYPE_GROUP)
			})

			Convey("message can be deleted by owner", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)
				err = rest.DeletePost(post.Id, ses.ClientId)
				So(err, ShouldBeNil)
				post2, err := rest.GetPost(post.Id, ses.ClientId)
				So(err, ShouldNotBeNil)
				So(post2, ShouldBeNil)
			})

			Convey("message should not have payload, if user does not allow", func() {
				h := http.Header{}
				h.Add("X-Forwarded-For", "208.72.139.54")
				post, err := rest.CreatePostWithHeader(groupChannel.Id, h, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)
				So(post.Payload, ShouldBeNil)
			})

			Convey("Message should have location if user allowed", func() {
				payload := make(map[string]interface{})
				payload["saveLocation"] = "Manisa"
				post, err := rest.CreatePostWithPayload(groupChannel.Id, payload, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)
				So(post.Payload, ShouldNotBeNil)
				So(*(post.Payload["saveLocation"]), ShouldEqual, "Manisa")
			})

			// handled by social worker
			Convey("message can be deleted by an admin", nil)
			Convey("message can not be edited by non-owner", nil)

			Convey("owner can post reply to message", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)

				reply, err := rest.AddReply(post.Id, groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(reply, ShouldNotBeNil)

				So(reply.AccountId, ShouldEqual, post.AccountId)

				cmc, err := rest.GetPostWithRelatedData(
					post.Id,
					&request.Query{
						AccountId: post.AccountId,
						GroupName: groupName,
					},
					ses.ClientId,
				)

				So(err, ShouldBeNil)
				So(cmc, ShouldNotBeNil)

				So(len(cmc.Replies), ShouldEqual, 1)

				So(cmc.Replies[0].Message.AccountId, ShouldEqual, post.AccountId)

			})

			Convey("we should be able to get only replies", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)

				reply, err := rest.AddReply(post.Id, groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(reply, ShouldNotBeNil)

				reply, err = rest.AddReply(post.Id, groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(reply, ShouldNotBeNil)

				replies, err := rest.GetReplies(post.Id, post.AccountId, groupName)
				So(err, ShouldBeNil)
				So(len(replies), ShouldEqual, 2)

			})

			Convey("we should be able to get replies with \"from\" query param", nil)

			Convey("non-owner can post reply to message", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)

				reply, err := rest.AddReply(post.Id, groupChannel.Id, nonOwnerSes.ClientId)
				So(err, ShouldBeNil)
				So(reply, ShouldNotBeNil)

				So(reply.AccountId, ShouldEqual, nonOwnerAccount.Id)

				cmc, err := rest.GetPostWithRelatedData(
					post.Id,
					&request.Query{
						AccountId: post.AccountId,
						GroupName: groupName,
					},
					ses.ClientId,
				)

				So(err, ShouldBeNil)
				So(cmc, ShouldNotBeNil)

				So(len(cmc.Replies), ShouldEqual, 1)

				So(cmc.Replies[0].Message.AccountId, ShouldEqual, nonOwnerAccount.Id)
			})

			Convey("reply can be deleted by owner", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)

				sesNonOwner, err := models.FetchOrCreateSession(nonOwnerAccount.Nick, groupName)
				So(err, ShouldBeNil)

				reply, err := rest.AddReply(post.Id, groupChannel.Id, nonOwnerSes.ClientId)
				So(err, ShouldBeNil)
				So(reply, ShouldNotBeNil)

				err = rest.DeletePost(reply.Id, sesNonOwner.ClientId)
				So(err, ShouldBeNil)

				cmc, err := rest.GetPostWithRelatedData(
					post.Id,
					&request.Query{
						AccountId: account.Id,
						GroupName: groupName,
					},
					ses.ClientId,
				)

				So(err, ShouldBeNil)
				So(cmc, ShouldNotBeNil)

				So(len(cmc.Replies), ShouldEqual, 0)

			})

			Convey("while deleting message, also replies should be deleted", func() {
				post, err := rest.CreatePost(groupChannel.Id, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)

				reply1, err := rest.AddReply(post.Id, groupChannel.Id, nonOwnerSes.ClientId)
				So(err, ShouldBeNil)
				So(reply1, ShouldNotBeNil)

				reply2, err := rest.AddReply(post.Id, groupChannel.Id, nonOwnerSes.ClientId)
				So(err, ShouldBeNil)
				So(reply2, ShouldNotBeNil)

				err = rest.DeletePost(post.Id, ses.ClientId)
				So(err, ShouldBeNil)

				cmc, err := rest.GetPostWithRelatedData(
					reply1.Id,
					&request.Query{
						AccountId: account.Id,
						GroupName: groupName,
					},
					ses.ClientId,
				)
				So(err, ShouldNotBeNil)
				So(cmc, ShouldBeNil)

				cmc, err = rest.GetPostWithRelatedData(
					reply2.Id,
					&request.Query{
						AccountId: account.Id,
						GroupName: groupName,
					},
					ses.ClientId,
				)

				So(err, ShouldNotBeNil)
				So(cmc, ShouldBeNil)

			})

			Convey("while deleting messages, they should be removed from all channels", nil)

			Convey("message can contain payload", func() {
				payload := make(map[string]interface{})
				payload["key1"] = "value1"
				payload["key2"] = 2
				payload["key3"] = true
				payload["key4"] = 3.4

				post, err := rest.CreatePostWithPayload(groupChannel.Id, payload, ses.ClientId)
				So(err, ShouldBeNil)
				So(post, ShouldNotBeNil)

				So(post.Payload, ShouldNotBeNil)
				So(*(post.Payload["key1"]), ShouldEqual, "value1")
				So(*(post.Payload["key2"]), ShouldEqual, "2")
				So(*(post.Payload["key3"]), ShouldEqual, "true")
				So(*(post.Payload["key4"]), ShouldEqual, "3.4")
			})
		})
	})
}
