import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/routes.dart';
import 'package:sourcepods/src/gravatar_component.dart';
import 'package:sourcepods/src/login/login_service.dart';
import 'package:sourcepods/src/user/user.dart';
import 'package:sourcepods/src/user/user_service.dart';

@Component(
  selector: 'sourcepods-header',
  templateUrl: 'header_component.html',
  styleUrls: const ['header_component.css'],
  directives: const [
    coreDirectives,
    routerDirectives,
    GravatarComponent,
  ],
  providers: const [
    UserService,
    LoginService,
  ],
)
class HeaderComponent implements OnInit {
  HeaderComponent(this._userService, this._loginService);

  final UserService _userService;
  final LoginService _loginService;

  String username;

  String email;

  @override
  void ngOnInit() {
    _userService.me().then((User user) {
      this.username = user.username;
      this.email = user.email;
    });
  }

  String usersUrl() => RoutePaths.userList.toUrl();

  String userProfile() =>
      RoutePaths.userProfile.toUrl(parameters: {'username': username});

  String loginUrl() => RoutePaths.login.toUrl();

  String settingsUrl() => RoutePaths.settings.toUrl();

  void logout() => this._loginService.logout();
}
