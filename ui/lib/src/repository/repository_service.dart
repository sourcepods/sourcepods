import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/api.dart';
import 'package:gitpods/src/api/api.dart' as api;
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/repository/repository_component.dart';
import 'package:gitpods/src/validation_exception.dart';
import 'package:http/http.dart';

@Injectable()
class RepositoryService {
  RepositoryService(this._http, this._api);

  final Client _http;
  final API _api;

  Future<RepositoryPage> get(String owner, String name) async {
    api.Repository r = await _api.repositories.getRepository(owner, name);
    return RepositoryPage(Repository.fromAPI(r));
  }

  Future<Repository> create(Repository repository) async {
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
