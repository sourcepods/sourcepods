import 'package:gitpods/repository.dart';

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

  factory User.fromJSON(Map<String, dynamic> data) {
    User user = new User(
      id: data['id'],
      email: data['email'],
      username: data['username'],
    );

    data['name'] != null ? user.name = data['name'] : '';

    if (data['created_at'] != null) {
      user.created =
      new DateTime.fromMillisecondsSinceEpoch(data['created_at'] * 1000);
    }

    if (data['updated_at'] != null) {
      user.updated =
      new DateTime.fromMillisecondsSinceEpoch(data['updated_at'] * 1000);
    }

    return user;
  }
}
