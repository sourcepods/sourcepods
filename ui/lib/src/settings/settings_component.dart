import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';

@Component(
  selector: 'settings',
  templateUrl: 'settings_component.html',
  directives: const [routerDirectives],
)
//@RouteConfig(const [
//  const Route(
//    path: '/profile',
//    name: 'Profile',
//    component: ProfileComponent,
//    useAsDefault: true,
//  ),
//  const Route(
//    path: '/account',
//    name: 'Account',
//    component: AccountComponent,
//  ),
//  const Route(
//    path: '/security',
//    name: 'Security',
//    component: SecurityComponent,
//  ),
//])
class SettingsComponent {}
