import 'package:gitpods/user.dart';

class Repository {
  String id;
  String name;
  String description;
  String website;
  String default_branch;
  bool private;
  bool bare;
  DateTime created;
  DateTime updated;

  int stars;
  int forks;
  Map<String, int> issueStats;
  Map<String, int> pullRequestStats;

  User owner;

  Repository({
    this.id,
    this.name,
    this.description,
    this.website,
    this.default_branch,
    this.private,
    this.bare,
    this.created,
    this.updated,
    this.stars,
    this.forks,
  });

  factory Repository.fromJSON(Map<String, dynamic> data) {
    Repository repository = new Repository(
      id: data['id'],
      name: data['name'],
    );

    data['description'] != null ? repository.description = data['description'] : '';
    data['website'] != null ? repository.website = data['website'] : '';
    data['default_branch'] != null ? repository.default_branch = data['default_branch'] : '';
    data['private'] != null ? repository.private = data['private'] : '';
    data['bare'] != null ? repository.bare = data['bare'] : '';
    data['stars'] != null ? repository.stars = data['stars'] : '';
    data['forks'] != null ? repository.forks = data['forks'] : '';

    if (data['issue_stats'] != null) {
      repository.issueStats = {
        'total': data['issue_stats']['total'],
        'open': data['issue_stats']['open'],
        'closed': data['issue_stats']['closed'],
      };
    }

    if (data['pull_request_stats'] != null) {
      repository.pullRequestStats = {
        'total': data['pull_request_stats']['total'],
        'open': data['pull_request_stats']['open'],
        'closed': data['pull_request_stats']['closed'],
      };
    }

    if (data['owner'] != null) {
      repository.owner = new User.fromJSON(data['owner']);
    }

    if (data['created_at'] != null) {
      repository.created =
      new DateTime.fromMillisecondsSinceEpoch(data['created_at'] * 1000);
    }

    if (data['updated_at'] != null) {
      repository.updated =
      new DateTime.fromMillisecondsSinceEpoch(data['updated_at'] * 1000);
    }

    return repository;
  }
}
