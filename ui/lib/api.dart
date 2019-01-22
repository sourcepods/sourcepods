import 'dart:html';
import 'dart:js_util';
import 'package:sourcepods/src/api/api.dart';

class API {
  API() {
    String url = getProperty(window, 'api');
    ApiClient client = new ApiClient(basePath: '$url/v1');

    this.repositories = RepositoriesApi(client);
    this.users = UsersApi(client);
  }

  RepositoriesApi repositories;
  UsersApi users;
}
