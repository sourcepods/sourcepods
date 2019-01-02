import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:http/http.dart';

@Injectable()
class UserService {
  UserService(this._http);

  final Client _http;

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

  Future<User> me() async {
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

  final usersQuery = '''
query UsersQuery {
  users {
    id
    email
    username
    name
    createdAt
    updatedAt
  }
}
''';

  Future<List<User>> list() async {
    var payload = json.encode({
      'query': usersQuery,
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = json.decode(resp.body);
    return body['data']['users']
        .map((json) => new User.fromJSON(json))
        .toList();
  }

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

  Future<UserProfile> profile(String username) async {
    var payload = json.encode({
      'query': userProfile,
      'variables': {
        'username': username,
      }
    });

    Response resp = await this._http.post('/api/query', body: payload);

    var body = json.decode(resp.body);

    User user = new User.fromJSON(body['data']['user']);
    List<Repository> repositories = body['data']['repositories']
        .map((json) => new Repository.fromJSON(json))
        .toList();

    return new UserProfile(
      user: user,
      repositories: repositories,
    );
  }

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

  Future<User> update(User user) async {
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
