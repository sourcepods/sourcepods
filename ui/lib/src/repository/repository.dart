import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/api/api.dart' as api;

class Repository {
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

  String id;
  String name;
  String description;
  String website;
  String defaultBranch;
  DateTime created;
  DateTime updated;

  User owner;
}
