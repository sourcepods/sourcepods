import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/gravatar_component.dart';
import 'package:gitpods/issues_component.dart';
import 'package:gitpods/login_component.dart';
import 'package:gitpods/login_service.dart';
import 'package:gitpods/pulls_component.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_list_component.dart';
import 'package:gitpods/user_profile_component.dart';
import 'package:gitpods/user_service.dart';
import 'package:gitpods/user_settings_component.dart';

@Component(
  selector: 'gitpods-app',
  templateUrl: 'app_component.html',
//  template: '<router-outlet></router-outlet>',
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES, Gravatar],
  providers: const [ROUTER_PROVIDERS, UserService, LoginService],
)
@RouteConfig(const[
  const Route(
    path: '/',
    name: 'UserList',
    component: UserListComponent,
    useAsDefault: true,
  ),
  const Route(
    path: '/login',
    name: 'Login',
    component: LoginComponent,
  ),
  const Route(
    path: '/issues',
    name: 'Issues',
    component: IssuesComponent,
  ),
  const Route(
    path: '/pulls',
    name: 'Pulls',
    component: PullsComponent,
  ),
  const Route(
    path: '/:username',
    name: 'UserProfile',
    component: UserProfileComponent,
  ),
  const Route(
    path: '/settings/profile',
    name: 'UserSettings',
    component: UserSettingsComponent,
  ),
])
class AppComponent implements OnInit {
  final Router _router;
  final UserService _userService;
  final LoginService _loginService;

  AppComponent(this._router, this._userService, this._loginService);

  User user;
  bool loading = false;

  @override
  ngOnInit() {
    this._userService.me()
        .then((User user) => this.user = user)
        .catchError((e) => this._router.navigate(['Login']));
  }

  logout() {
    this._loginService.logout();
  }
}
