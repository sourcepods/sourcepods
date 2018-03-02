import 'dart:html';
import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';

@Component(
  selector: 'account',
  templateUrl: 'account_component.html',
  directives: const [CORE_DIRECTIVES, formDirectives],
)
class AccountComponent {
  bool loading;
  String passwordCurrent;
  String password1;
  String password2;

  void changePassword(Event e) {
    e.preventDefault();
    loading = true;
    print('$passwordCurrent, $password1, $password2, $loading');
  }
}
