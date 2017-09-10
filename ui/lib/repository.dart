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
