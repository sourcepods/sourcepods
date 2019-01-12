import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/api/api.dart' as api;

class Repository {
  String id;
  String name;
  String description;
  String website;
  String defaultBranch;
  DateTime created;
  DateTime updated;

  User owner;
  List<RepositoryBranch> branches;

  Repository({
    this.id,
    this.name,
    this.description,
    this.website,
    this.defaultBranch,
    this.created,
    this.updated,
  });

  factory Repository.fromAPI(api.Repository r) {
    return Repository(
      id: r.id,
      name: r.name,
      description: r.description,
      website: r.website,
      defaultBranch: r.defaultBranch,
      created: r.createdAt,
      updated: r.updatedAt,
    );
  }


  factory Repository.fromJSON(Map<String, dynamic> data) {
    Repository repository = new Repository(
      id: data['id'],
      name: data['name'],
    );

    data['description'] != null ? repository.description = data['description'] : '';
    data['website'] != null ? repository.website = data['website'] : '';
    data['defaultBranch'] != null ? repository.defaultBranch = data['defaultBranch'] : '';

    if (data['owner'] != null) {
      repository.owner = new User.fromJSON(data['owner']);
    }

    if (data['createdAt'] != null) {
      repository.created = DateTime.parse(data['createdAt']);
    }

    if (data['updatedAt'] != null) {
      repository.updated = DateTime.parse(data['updatedAt']);
    }

    if (data['branches'] != null) {
      repository.branches = data['branches']
          .map((json) => new RepositoryBranch.fromJSON(json))
          .toList();
    }

    return repository;
  }

}

class RepositoryBranch {
  String name;
  String sha1;
  String type;
  bool protected;

  RepositoryBranch({
    this.name,
    this.sha1,
    this.type,
    this.protected,
  });

  factory RepositoryBranch.fromJSON(Map<String, dynamic> data) {
    return new RepositoryBranch(
      name: data['name'],
      sha1: data['sha1'],
      type: data['type'],
      protected: data['protected'],
    );
  }
}
