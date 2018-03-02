import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/settings/account/account_component.dart';
import 'package:gitpods/src/settings/profile/profile_component.dart';
import 'package:gitpods/src/settings/security/security_component.dart';

@Component(
  selector: 'settings',
  templateUrl: 'settings_component.html',
  directives: const [ROUTER_DIRECTIVES],
)
@RouteConfig(const [
  const Route(
    path: '/profile',
    name: 'Profile',
    component: ProfileComponent,
    useAsDefault: true,
  ),
  const Route(
    path: '/account',
    name: 'Account',
    component: AccountComponent,
  ),
  const Route(
    path: '/security',
    name: 'Security',
    component: SecurityComponent,
  ),
])
class SettingsComponent {}
