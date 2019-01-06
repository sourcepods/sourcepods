import 'package:gitpods/src/api/api.dart';

class API {
  API() {
    // TODO: What's the best way to inject the correct basePath from the backend here?
    ApiClient client = new ApiClient(basePath: 'http://localhost:3000/api/v1');

    this.users = UsersApi(client);
  }

  UsersApi users;
}
