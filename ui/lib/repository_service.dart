import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:gitpods/repository.dart';
import 'package:gitpods/validation_exception.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';

const repositoryGet = '''
query (\$owner: String!, \$name: String!) {
  repository(owner: \$owner, name: \$name) {
    id
    name
    description
    website
    default_branch
    private
    created_at
    updated_at
    stars
    forks
    issue_stats {
      total
      open
      closed
    }
    pull_request_stats {
      total
      open
      closed
    }
  }
  tree(owner: \$owner, name: \$name) {
    type
    file
    commit {
      hash
      subject
      author {
        name
        email
        date
      }
      committer {
        name
        email
        date
      }
    }
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

  Future<Repository> get(String owner, String name) async {
    var payload = JSON.encode({
      'query': repositoryGet,
      'variables': {
        'owner': owner,
        'name': name,
      },
    });

    Response resp = await this._http.post('/api/query', body: payload);
    var body = JSON.decode(resp.body);

    if (body['errors'] != null) {
      throw new Exception(body['errors'][0]['message']);
    }

    return new Repository.fromJSON(body['data']['repository']);
  }

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
