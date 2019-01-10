import 'package:gitpods/src/api/api.dart' as api;
import 'package:gitpods/src/repository/repository.dart';

class User {
  String id;
  String email;
  String username;
  String name;
  DateTime created;
  DateTime updated;

  List<Repository> repositories;

  User({
    this.id,
    this.email,
    this.username,
    this.name,
    this.created,
    this.updated,
  });

  factory User.fromAPI(api.User user) {
    return User(
      id: user.id,
      email: user.email,
      username: user.username,
      name: user.name,
      created: user.createdAt,
      updated: user.updatedAt,
    );
  }

  factory User.fromJSON(Map<String, dynamic> data) {
    User user = new User(
      id: data['id'],
      email: data['email'],
      username: data['username'],
    );

    data['name'] != null ? user.name = data['name'] : '';

    if (data['repositories'] != null) {
      user.repositories = data['repositories']
          .map((data) => new Repository.fromJSON(data))
          .toList();
    }

    if (data['createdAt'] != null) {
      user.created = DateTime.parse(data['createdAt']);
    }

    if (data['updatedAt'] != null) {
      user.updated = DateTime.parse(data['updatedAt']);
    }

    return user;
  }
}
