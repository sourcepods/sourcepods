package store

//// TODO: Get rid of this store, its data is accessible via other stores combined
//
//type UserInMemory struct {
//	userStore         *UsersInMemory
//	repositoriesStore *RepositoriesInMemory
//}
//
//func NewUserInMemory(users *UsersInMemory, repositories *RepositoriesInMemory) *UserInMemory {
//	return &UserInMemory{userStore: users, repositoriesStore: repositories}
//}
//
//// TODO: Get this form the UsersStore
//func (s *UserInMemory) GetUser(username string) (*gitpods.User, error) {
//	return s.userStore.GetUser(username)
//}
//
//// TODO: Get this from the UsersRepositoriesStore
//func (s *UserInMemory) GetUserRepositories(username string) (*gitpods.User, []*gitpods.Repository, error) {
//	user, err := s.GetUser(username)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	repositories, err := s.repositoriesStore.List()
//	if err != nil {
//		return nil, nil, err
//	}
//
//	var userRepos []*gitpods.Repository
//	for _, repo := range repositories {
//		if repo.Owner == user {
//			userRepos = append(userRepos, &gitpods.Repository{
//				ID:            repo.ID,
//				Name:          repo.Name,
//				Description:   repo.Description,
//				Website:       repo.Website,
//				DefaultBranch: repo.DefaultBranch,
//				Private:       repo.Private,
//				Bare:          repo.Bare,
//				Created:       repo.Created,
//				Updated:       repo.Updated,
//
//				OwnerID: repo.Owner.ID,
//			})
//		}
//	}
//
//	return user, userRepos, nil
//}
