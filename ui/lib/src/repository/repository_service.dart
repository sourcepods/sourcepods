import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/repository/repository_component.dart';
import 'package:gitpods/src/validation_exception.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';

const repositoryGet = '''
query (\$owner: String!, \$name: String!) {
  repository(owner: \$owner, name: \$name) {
    id
    name
    description
    website
  }
}
''';

const repositoryCreate = '''
mutation (\$repository: newRepository!) {
  createRepository(repository: \$repository) {
    id
    name
    description
    website
    createdAt
    updatedAt
    owner {
      id
      username
      name
      email
    }
  }
}
''';

@Injectable()
class RepositoryService {
  final BrowserClient _http;

  RepositoryService(this._http);

  Future<RepositoryPage> get(String owner, String name) async {
    var payload = json.encode({
      'query': repositoryGet,
      'variables': {
        'owner': owner,
        'name': name,
      },
    });

    Response resp = await this._http.post('/api/query', body: payload);
    var body = json.decode(resp.body);

    if (body['errors'] != null) {
      throw new Exception(body['errors'][0]['message']);
    }

    Repository repository = new Repository.fromJSON(body['data']['repository']);
    return new RepositoryPage(repository);
  }

  Future<Repository> create(Repository repository) async {
    var payload = json.encode({
      'query': repositoryCreate,
      'variables': {
        'repository': {
          'name': repository.name,
          'description': repository.description,
          'website': repository.website,
        },
      },
    });

    Response resp = await this._http.post('/api/query', body: payload);
    var body = json.decode(resp.body);

    if (body['errors'] != null) {
      throw new ValidationException(body['errors'][0]['message']);
    }

    return new Repository.fromJSON(body['data']['createRepository']);
  }
}
