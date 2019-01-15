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
}
