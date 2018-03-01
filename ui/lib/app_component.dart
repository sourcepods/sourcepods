import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/login/login_component.dart';
import 'package:gitpods/src/login/login_service.dart';
import 'package:gitpods/src/repository/repository_component.dart';
import 'package:gitpods/src/repository/repository_create_component.dart';
import 'package:gitpods/src/not_found_component.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/user/user_list_component.dart';
import 'package:gitpods/src/user/user_profile_component.dart';
import 'package:gitpods/src/user/user_service.dart';
import 'package:gitpods/src/settings/settings_component.dart';

@Component(
  selector: 'gitpods-app',
  templateUrl: 'app_component.html',
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES, GravatarComponent],
  providers: const [ROUTER_PROVIDERS, UserService, LoginService],
)
@RouteConfig(const [
  const Redirect(
    path: '/',
    redirectTo: const ['UserList'],
  ),
  const Route(
    path: '/users',
    name: 'UserList',
    component: UserListComponent,
  ),
  const Route(
    path: '/login',
    name: 'Login',
    component: LoginComponent,
  ),
  const Route(
    path: '/:username',
    name: 'UserProfile',
    component: UserProfileComponent,
  ),
  const Route(
    path: '/settings/...',
    name: 'Settings',
    component: SettingsComponent,
  ),
  const Route(
    path: '/:owner/:name',
    name: 'Repository',
    component: RepositoryComponent,
  ),
  const Route(
    path: '/new',
    name: 'RepositoryCreate',
    component: RepositoryCreateComponent,
  ),
  const Route(
    path: '/**',
    name: 'NotFound',
    component: NotFoundComponent,
  )
])
class AppComponent implements OnInit {
  final Router _router;
  final UserService _userService;
  final LoginService _loginService;

  AppComponent(this._router, this._userService, this._loginService);

  User user;
  bool loading = false;

  @override
  void ngOnInit() {
    this._userService.me()
        .then((User user) => this.user = user)
        .catchError((e) => this._router.navigate(['Login']));
  }

  void logout() {
    this._loginService.logout();
  }
}
