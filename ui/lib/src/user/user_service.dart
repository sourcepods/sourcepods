import 'dart:async';

import 'package:angular/angular.dart';
import 'package:gitpods/api.dart';
import 'package:gitpods/src/api/api.dart' as api;
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/user/user.dart';

@Injectable()
class UserService {
  UserService(this._api);

  final API _api;

  Future<User> me() async {
    api.User apiUser = await _api.users.getUserMe();
    return User.fromAPI(apiUser);
  }

  Future<List<User>> list() async {
    List<api.User> list = await _api.users.listUsers();
    return list.map((api.User apiUser) => User.fromAPI(apiUser)).toList();
  }

  Future<UserProfile> profile(String username) async {
    List responses = await Future.wait([
      _api.users.getUser(username),
      _api.repositories.getOwnerRepositories(username),
    ]);

    api.User user = responses[0];
    List<api.Repository> repositories = responses[1];

    return UserProfile(
      user: User.fromAPI(user),
      repositories: (repositories)
          .map((api.Repository r) => Repository.fromAPI(r))
          .toList(),
    );
  }

  Future<User> update(User user) async {
    api.UpdatedUser updated = new api.UpdatedUser();
    updated.name = user.name;

    api.User apiUser = await _api.users.updateUser(user.username, updated);

    return User(
      id: apiUser.id,
      email: apiUser.email,
      username: apiUser.username,
      name: apiUser.name,
      created: apiUser.createdAt,
      updated: apiUser.updatedAt,
    );
  }
}

class UserProfile {
  UserProfile({this.repositories, this.user});

  List<Repository> repositories;
  User user;
}

class UnauthorizedError extends Error {
  UnauthorizedError(this.message);

  final String message;

  @override
  String toString() => "Unauthorizted state: $message";
}
