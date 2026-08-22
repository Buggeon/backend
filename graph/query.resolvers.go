package graph

import (
	"context"
	"fmt"
	"time"

	"bugtracker/graph/model"
	"bugtracker/internal/models"
)

func (r *queryResolver) Project(ctx context.Context, projectID string) (*model.Project, error) {
	project, err := r.ProjectService.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	return r.toGraphQLProject(context.Background(), &project), nil
}

func (r *queryResolver) Projects(ctx context.Context, userID string) ([]*model.Project, error) {

	projects, err := r.ProjectService.GetProjects(userID)

	if err != nil {
		return nil, err
	}

	result := make([]*model.Project, len(projects))

	for i, p := range projects {

		result[i] = r.toGraphQLProject(context.Background(), &p)
	}

	fmt.Println(result)

	return result, nil
}

func (r *Resolver) toGraphQLProject(ctx context.Context, p *models.Project) *model.Project {

	if p == nil {
		return nil
	}

	var lead *models.Member

	member, err := r.MemberService.GetMember(p.LeadID.Hex())

	if err == nil {
		lead = &member
	}

	return &model.Project{
		ID:          p.ID.Hex(),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		Name:        p.Name,
		Description: p.Description,
		LogoURL:     p.LogoUrl,
		Progress:    int32(p.Progress),
		Lead:        r.toGraphQLMember(context.Background(), lead),
	}
}

func (r *Resolver) toGraphQLMember(ctx context.Context, m *models.Member) *model.Member {

	if m == nil {
		return nil
	}

	fmt.Println(m.UserID)

	user, err := r.UserService.GetUser(m.UserID.Hex())

	fmt.Println("---USER MODEL---")
	fmt.Println(user)

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
