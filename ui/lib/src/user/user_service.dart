import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/src/api/api.dart' as api;
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:http/http.dart';

@Injectable()
class UserService {
  UserService(this._http);

  final Client _http;
  final api.UsersApi _usersApi = new api.UsersApi();

  Future<User> me() async {
    final userMe = '''
query me {
  me {
    id
    email
    username
    name
    createdAt
    updatedAt
  }
}
''';

    String payload = json.encode({
      'query': userMe,
    });

    Response resp = await _http.post('/api/query', body: payload);

    if (resp.statusCode != 200) {
      var body = json.decode(resp.body);
      throw new UnauthorizedError(body['errors'][0]['detail']);
    }

    var body = json.decode(resp.body);
    return new User.fromJSON(body['data']['me']);
  }

  Future<List<User>> list() async {
    List<api.User> list = await _usersApi.listUsers();

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
    final userUpdate = '''
mutation updateUser(\$id: ID!, \$user: UpdatedUser!) {
  user: updateUser(id: \$id, user: \$user) {
    id
    email
    username
    name
    createdAt
    updatedAt
  }
}
''';

    var payload = json.encode({
      'query': userUpdate,
      'variables': {
        'id': user.id,
        'user': {
          'name': user.name,
        },
      },
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = json.decode(resp.body);
    return new User.fromJSON(body['data']['user']);
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
