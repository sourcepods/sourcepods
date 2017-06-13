package resolver

// Resolver is the root with all other resolvers embedded which really implement all funcs.
type Resolver struct {
	*UserResolver
	*RepositoryResolver
}

var (
	// Schema to build the GraphQL API against.
	Schema = `
	schema {
		query: Query
		mutation: Mutation
	}` + Query + Mutation

	Query = `
	# The query type, represents all of the entry points into our object graph
	type Query {
		me: User
		user(id: ID, username: String): User
		users: [User]!
		repository(id: ID, owner: String, name: String): Repository
		repositories(owner: String!): [Repository]!
	}
	type User {
		id: ID!
		email: String!
		username: String!
		name: String!
		created_at: Int!
		updated_at: Int!
		repositories: [Repository]!
	}
	# Something about a repository
	type Repository {
		id: ID!
		name: String!
		description: String!
		website: String!
		default_branch: String!
		private: Boolean!
		bare: Boolean!
		created_at: Int!
		updated_at: Int!
		stars: Int!
		forks: Int!
		issue_stats: IssueStats!
		pull_request_stats: PullRequestStats!
	}
	interface OpenClosedStats {
		total: Int!
		open: Int!
		closed: Int!
	}
	type IssueStats implements OpenClosedStats {
		total: Int!
		open: Int!
		closed: Int!
	}
	type PullRequestStats implements OpenClosedStats {
		total: Int!
		open: Int!
		closed: Int!
	}`
	Mutation = `
	# The mutation type, represents all updates we can make to our data
	type Mutation {
		updateUser(id: ID!, user: updatedUser!): User!
	}
	input updatedUser {
		name: String!
	}`
)
