import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_service.dart';

@Component(
  selector: 'gitpods-navbar',
  templateUrl: 'navbar_component.html',
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES],
  providers: const [ROUTER_PROVIDERS, UserService],
)
class NavbarComponent implements OnInit {
  final UserService _userService;
  final Router _router;

  NavbarComponent(this._userService, this._router);

  bool loading = false;
  User user;

  @override
  ngOnInit() {
    this._userService.me()
        .then((User user) => this.user = user)
        .catchError((e) => this._router.navigate(['Login']));
  }
}
