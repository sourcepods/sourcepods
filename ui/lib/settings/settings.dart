import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/settings/profile.dart';
import 'package:gitpods/settings/security.dart';

@Component(
  selector: 'settings',
  templateUrl: 'settings.html',
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
    path: '/security',
    name: 'Security',
    component: SecurityComponent,
  ),
])
class SettingsComponent {}
