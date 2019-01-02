import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/settings/account/account_component.template.dart';
import 'package:gitpods/src/settings/profile/profile_component.template.dart';
import 'package:gitpods/src/settings/security/security_component.template.dart';

final profileRoute = new RoutePath(path: '/', useAsDefault: true);
final accountRoute = new RoutePath(path: '/account');
final securityRoute = new RoutePath(path: '/security');

@Injectable()
class Routes {
  final List<RouteDefinition> all = [
    new RouteDefinition(
      routePath: profileRoute,
      component: ProfileComponentNgFactory,
    ),
    new RouteDefinition(
      routePath: accountRoute,
      component: AccountComponentNgFactory,
    ),
    new RouteDefinition(
      routePath: securityRoute,
      component: SecurityComponentNgFactory,
    )
  ];
}
