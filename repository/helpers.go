package repository

import "github.com/anilaydinn/socium-be/model"

func convertUserModelToUserEntity(user model.User) UserEntity {
	return UserEntity{
		ID:                   user.ID,
		Name:                 user.Name,
		Surname:              user.Surname,
		Email:                user.Email,
		BirthDate:            user.BirthDate,
		Description:          user.Description,
		ProfileImage:         user.ProfileImage,
		FriendRequestUserIDs: user.FriendRequestUserIDs,
		FriendIDs:            user.FriendIDs,
		Password:             user.Password,
		UserType:             user.UserType,
		IsActivated:          user.IsActivated,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
		Latitude:             user.Latitude,
		Longitude:            user.Longitude,
	}
}

func convertUserEntityToUserModel(userEntity UserEntity) model.User {
	return model.User{
		ID:                   userEntity.ID,
		Name:                 userEntity.Name,
		Surname:              userEntity.Surname,
		Email:                userEntity.Email,
		BirthDate:            userEntity.BirthDate,
		Description:          userEntity.Description,
		ProfileImage:         userEntity.ProfileImage,
		FriendRequestUserIDs: userEntity.FriendRequestUserIDs,
		FriendIDs:            userEntity.FriendIDs,
		Password:             userEntity.Password,
		UserType:             userEntity.UserType,
		IsActivated:          userEntity.IsActivated,
		CreatedAt:            userEntity.CreatedAt,
		UpdatedAt:            userEntity.UpdatedAt,
		Latitude:             userEntity.Latitude,
		Longitude:            userEntity.Longitude,
	}
}

func convertPostModelToPostEntity(post model.Post) PostEntity {
	return PostEntity{
		ID:              post.ID,
		UserID:          post.UserID,
		Description:     post.Description,
		Image:           post.Image,
		IsPrivate:       post.IsPrivate,
		WhoLikesUserIDs: post.WhoLikesUserIDs,
		CommentIDs:      post.CommentIDs,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	}
}

func convertPostEntityToPostModel(postEntity PostEntity) model.Post {
	return model.Post{
		ID:              postEntity.ID,
		UserID:          postEntity.UserID,
		Description:     postEntity.Description,
		Image:           postEntity.Image,
		IsPrivate:       postEntity.IsPrivate,
		WhoLikesUserIDs: postEntity.WhoLikesUserIDs,
		CommentIDs:      postEntity.CommentIDs,
		CreatedAt:       postEntity.CreatedAt,
		UpdatedAt:       postEntity.UpdatedAt,
	}
}

func convertCommentModelToCommentEntity(comment model.Comment) CommentEntity {
	return CommentEntity{
		ID:        comment.ID,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func convertCommentEntityToCommentModel(commentEntity CommentEntity) model.Comment {
	return model.Comment{
		ID:        commentEntity.ID,
		UserID:    commentEntity.UserID,
		PostID:    commentEntity.PostID,
		Content:   commentEntity.Content,
		CreatedAt: commentEntity.CreatedAt,
		UpdatedAt: commentEntity.UpdatedAt,
	}
}

func convertContactModelToContactEntity(contact model.Contact) ContactEntity {
	return ContactEntity{
		ID:      contact.ID,
		Name:    contact.Name,
		Surname: contact.Surname,
		Email:   contact.Email,
		Message: contact.Message,
	}
}

func convertContactEntityToContactModel(contactEntity ContactEntity) model.Contact {
	return model.Contact{
		ID:      contactEntity.ID,
		Name:    contactEntity.Name,
		Surname: contactEntity.Surname,
		Email:   contactEntity.Email,
		Message: contactEntity.Message,
	}
}
