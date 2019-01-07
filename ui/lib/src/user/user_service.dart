import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/api.dart';
import 'package:gitpods/src/api/api.dart' as api;
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:http/http.dart';

@Injectable()
class UserService {
  UserService(this._http, this._api);

  final Client _http;
  final API _api;

  Future<User> me() async {
    api.User apiUser = await _api.users.getUserMe();

    return User(
      id: apiUser.id,
      email: apiUser.email,
      username: apiUser.username,
      name: apiUser.name,
      created: apiUser.createdAt,
      updated: apiUser.updatedAt,
    );
  }

  Future<List<User>> list() async {
    List<api.User> list = await _api.users.listUsers();

    List<User> users = list.map((api.User user) {
      return User(
        id: user.id,
        email: user.email,
        username: user.username,
        name: user.name,
        created: user.createdAt,
        updated: user.updatedAt,
      );
    }).toList();

    return users;
  }

  Future<UserProfile> profile(String username) async {
    final userProfile = '''
query userProfile(\$username: String!) {
  user(username: \$username) {
    id
    username
    name
    email
    createdAt
    updatedAt
  }
  repositories(owner: \$username) {
    id
    name
    description
  }
}
''';

    var payload = json.encode({
      'query': userProfile,
      'variables': {
        'username': username,
      }
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = json.decode(resp.body);

    User user = new User.fromJSON(body['data']['user']);

    List<Repository> repositories = (body['data']['repositories'] as List)
        .map((json) => new Repository.fromJSON(json))
        .toList();

    return new UserProfile(
      user: user,
      repositories: repositories,
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
