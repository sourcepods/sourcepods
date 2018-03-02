import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/login/login_service.dart';

@Component(
  selector: 'gitpods-header',
  templateUrl: 'header_component.html',
  styleUrls: const ['header_component.css'],
  directives: const [
    COMMON_DIRECTIVES,
    ROUTER_DIRECTIVES,
    GravatarComponent,
  ],
)
class HeaderComponent {
  final LoginService _loginService;

  @Input('username')
  String username;

  @Input('email')
  String email;

  HeaderComponent(this._loginService);

  void logout() {
    this._loginService.logout();
  }
}
