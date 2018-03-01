import 'dart:async';
import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:gitpods/src/login/login_service.dart';

@Component(
  selector: 'gitpods-selector',
  templateUrl: 'login_component.html',
  styleUrls: const ['login_component.css'],
  directives: const[COMMON_DIRECTIVES, formDirectives],
  providers: const[LoginService],
)
class LoginComponent {
  final LoginService _loginService;

  String email;
  String password;

  String error = '';

  LoginComponent(this._loginService);

  void submit(Event event) {
    event.preventDefault();
    this.error = '';

    Future<String> resp = this._loginService.login(this.email, this.password);
    resp.then((String error) => this.error = error);
  }
}
