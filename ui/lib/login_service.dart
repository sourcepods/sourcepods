import 'dart:async';
import 'dart:convert';
import 'dart:html';

import 'package:angular/angular.dart';
import 'package:http/browser_client.dart';

@Injectable()
class LoginService {
  BrowserClient _http;

  LoginService(this._http);

  Future<String> login(String email, String password) async {
    try {
      Map body = {'email': email, 'password': password};
      final response =
          await _http.post('/api/authorize', body: JSON.encode(body));

      if (response.statusCode == 200) {
        // Reload the page to have the new cookie set.
        window.location.assign('/');
        return '';
      }

      var data = JSON.decode(response.body);

      return data['errors'][0]['detail'];
    } catch (e) {
      return e.toString();
    }
  }

  void logout() {
    window.location.assign('/api/sessions/logout');
  }
}
