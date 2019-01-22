import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/routes.dart' as _parent;
import 'package:sourcepods/src/settings/account/account_component.template.dart';
import 'package:sourcepods/src/settings/profile/profile_component.template.dart';
import 'package:sourcepods/src/settings/security/security_component.template.dart';

class RoutesPaths {
  static final account = RoutePath(
    path: '/account',
    parent: _parent.RoutePaths.settings,
  );
  static final security = RoutePath(
    path: '/security',
    parent: _parent.RoutePaths.settings,
  );
  static final profile = RoutePath(
    path: '',
    parent: _parent.RoutePaths.settings,
    useAsDefault: true,
  );
}

class Routes {
  static final List<RouteDefinition> all = [
    new RouteDefinition(
      routePath: RoutesPaths.account,
      component: AccountComponentNgFactory,
    ),
    new RouteDefinition(
      routePath: RoutesPaths.security,
      component: SecurityComponentNgFactory,
    ),
    new RouteDefinition(
      routePath: RoutesPaths.profile,
      component: ProfileComponentNgFactory,
      useAsDefault: true,
    ),
  ];
}
