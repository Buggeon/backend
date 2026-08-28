package graph

import (
	"context"
	"fmt"
	"time"

	"bugtracker/graph/model"
	"bugtracker/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Project(ctx context.Context, projectID string) (*model.Project, error) {

	project, err := r.ProjectService.GetProject(projectID)

	if err != nil {
		return nil, err
	}

	return r.toGraphQLProject(context.Background(), &project), nil
}

func (r *queryResolver) Projects(ctx context.Context, userID string) ([]*model.Project, error) {

	fmt.Println("===START TO LOAD PROJECTS===")
	fmt.Println(userID)

	projects, err := r.ProjectService.GetProjects(userID)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Project, len(projects))

	for i, p := range projects {

		result[i] = r.toGraphQLProject(context.Background(), &p)
	}

	return result, nil
}

func (r *queryResolver) User(ctx context.Context, userID string) (*model.User, error) {

	user, err := r.UserService.GetUser(userID)

	if err != nil {
		return &model.User{}, err
	}

	return r.toGraphQLUser(context.Background(), user), nil

}

func (r *queryResolver) Board(ctx context.Context, boardID string) (*model.Board, error) {

	board, err := r.BoardService.GetBoard(boardID)

	if err != nil {
		return &model.Board{}, err
	}

	return r.toGraphQLBoard(context.Background(), &board), nil

}

func (r *queryResolver) Boards(ctx context.Context, projectID string) ([]*model.Board, error) {

	boards, err := r.BoardService.GetBoards(projectID)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Board, len(boards))

	for i, b := range boards {

		result[i] = r.toGraphQLBoard(context.Background(), &b)

	}

	return result, nil

}

func (r *queryResolver) Card(ctx context.Context, cardID string) (*model.Card, error) {

	card, err := r.CardService.GetCard(cardID)

	if err != nil {
		return &model.Card{}, err
	}

	return r.toGraphQLCard(context.Background(), &card), nil

}

func (r *queryResolver) Cards(ctx context.Context, boardID string) ([]*model.Card, error) {

	cards, err := r.CardService.GetCards(boardID)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Card, len(cards))

	for i, c := range cards {

		result[i] = r.toGraphQLCard(context.Background(), &c)

	}

	return result, nil

}

func (r *queryResolver) Message(ctx context.Context, messageID string) (*model.Message, error) {

	message, err := r.MessageService.GetMessage(messageID)

	if err != nil {
		return &model.Message{}, err
	}

	return r.toGraphQLMessage(context.Background(), &message), nil

}

func (r *queryResolver) Messages(ctx context.Context, cardID string) ([]*model.Message, error) {

	messages, err := r.MessageService.GetMessages(cardID)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Message, len(messages))

	for i, m := range messages {

		result[i] = r.toGraphQLMessage(context.Background(), &m)

	}

	return result, nil

}

func (r *queryResolver) Member(ctx context.Context, memberID string) (*model.Member, error) {

	member, err := r.MemberService.GetMember(memberID)

	if err != nil {
		return &model.Member{}, err
	}

	return r.toGraphQLMember(context.Background(), &member), nil

}

func (r *queryResolver) Members(ctx context.Context, projectID string) ([]*model.Member, error) {

	members, err := r.MemberService.GetMembers(projectID)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Member, len(members))

	for i, m := range members {

		result[i] = r.toGraphQLMember(context.Background(), &m)

	}

	return result, nil

}

func (r *Resolver) toGraphQLProject(ctx context.Context, p *models.Project) *model.Project {

	if p == nil {
		return nil
	}

	lead, err := r.MemberService.GetMember(p.LeadID.Hex())

	if err != nil {
		return &model.Project{}
	}

	var members []*model.Member
	var boards []*model.Board

	for _, memberID := range p.Members {

		member, err := r.MemberService.GetMember(memberID.Hex())

		if err != nil {
			return &model.Project{}
		}

		graphMember := r.toGraphQLMember(context.Background(), &member)

		members = append(members, graphMember)

	}

	for _, boardID := range p.Boards {

		board, err := r.BoardService.GetBoard(boardID.Hex())

		if err != nil {
			return &model.Project{}
		}

		graphBoard := r.toGraphQLBoard(context.Background(), &board)

		boards = append(boards, graphBoard)

	}

	return &model.Project{
		ID:          p.ID.Hex(),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		Name:        p.Name,
		Description: p.Description,
		LogoURL:     p.LogoUrl,
		Progress:    int32(p.Progress),
		Lead:        r.toGraphQLMember(context.Background(), &lead),
		Members:     members,
		Boards:      boards,
	}
}

func (r *Resolver) toGraphQLMember(ctx context.Context, m *models.Member) *model.Member {

	if m == nil {
		return nil
	}

	user, err := r.UserService.GetUser(m.UserID.Hex())

	if err != nil {
		return nil
	}

	return &model.Member{
		ID:         m.ID.Hex(),
		Role:       m.Role,
		Directions: m.Directions,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
		ProjectID:  m.ProjectID.Hex(),
		User:       r.toGraphQLUser(context.Background(), user),
	}

}

func (r *Resolver) toGraphQLUser(ctx context.Context, u *models.User) *model.User {

	if u == nil {
		return nil
	}

	return &model.User{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Login:     u.Login,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func (r *Resolver) toGraphQLBoard(ctx context.Context, b *models.Board) *model.Board {

	if b == nil {
		return nil
	}

	var cards []*model.Card

	for _, cardID := range b.Cards {

		card, err := r.CardService.GetCard(cardID.Hex())

		if err != nil {
			return &model.Board{}
		}

		graphCard := r.toGraphQLCard(context.TODO(), &card)

		cards = append(cards, graphCard)

	}

	return &model.Board{
		ID:        b.ID.Hex(),
		Name:      b.Name,
		Direction: b.Direction,
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
		Cards:     cards,
	}
}

func (r *Resolver) toGraphQLCard(ctx context.Context, c *models.Card) *model.Card {

	if c == nil {
		return nil
	}

	var assignees []*model.Member
	var messages []*model.Message

	for _, memberID := range c.Assignees {

		member, err := r.MemberService.GetMember(memberID.Hex())

		if err != nil {
			return &model.Card{}
		}

		graphMember := r.toGraphQLMember(context.Background(), &member)

		assignees = append(assignees, graphMember)

	}

	for _, messageID := range c.Messages {

		message, err := r.MessageService.GetMessage(messageID.Hex())

		if err != nil {
			return &model.Card{}
		}

		graphMessage := r.toGraphQLMessage(context.Background(), &message)

		messages = append(messages, graphMessage)

	}

	return &model.Card{
		ID:        c.ID.Hex(),
		Title:     c.Title,
		Content:   c.Content,
		Priority:  int32(c.Priority),
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		BoardID:   c.BoardID.Hex(),
		Assigness: assignees,
		Messages:  messages,
	}
}

func (r *Resolver) toGraphQLMessage(ctx context.Context, m *models.Message) *model.Message {
	if m == nil {
		return nil
	}

	sender, err := r.MemberService.GetMember(m.SenderID.Hex())
	if err != nil {
		return &model.Message{
			ID:        m.ID.Hex(),
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
			UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		}
	}

	var replyTo *model.Message
	if m.ReplyTo != primitive.NilObjectID {
		replyMessage, err := r.MessageService.GetMessage(m.ReplyTo.Hex())
		if err == nil {
			replyTo = r.toGraphQLMessage(context.Background(), &replyMessage)
		}
	}

	return &model.Message{
		ID:        m.ID.Hex(),
		Sender:    r.toGraphQLMember(context.Background(), &sender),
		CardID:    m.CardID.Hex(),
		ReplyTo:   replyTo,
		Content:   m.Content,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}
