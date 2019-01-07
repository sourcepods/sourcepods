import 'dart:html';
import 'dart:js_util';
import 'package:gitpods/src/api/api.dart';

class API {
  API() {
    String url = getProperty(window, 'api');
    ApiClient client = new ApiClient(basePath: '$url/v1');

    this.users = UsersApi(client);
  }

  UsersApi users;
}
