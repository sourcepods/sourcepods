import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/repository.dart';
import 'package:gitpods/validation_exception.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';

const repositoryCreate = '''
mutation (\$repository: newRepository!) {
  createRepository(repository: \$repository) {
    id
    name
    description
    website
    created_at
    updated_at
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

  Future<Repository> create(Repository repository) async {
    var payload = JSON.encode({
      'query': repositoryCreate,
      'variables': {
        'repository': {
          'name': repository.name,
          'description': repository.description,
          'website': repository.website,
          'private': repository.private,
        },
      },
    });

    Response resp = await this._http.post('/api/query', body: payload);
    var body = JSON.decode(resp.body);

    if (body['errors'] != null) {
      throw new ValidationException(body['errors'][0]['message']);
    }

    return new Repository.fromJSON(body['data']['createRepository']);
  }
}
