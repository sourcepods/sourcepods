import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/repository/commits/commits_component.template.dart';
import 'package:gitpods/src/repository/files/files_component.template.dart';
import 'package:gitpods/src/repository/settings/settings_component.template.dart';
import 'package:gitpods/routes.dart' as global;

class RoutePaths {
  static final files = RoutePath(
    path: '/',
    parent: global.RoutePaths.repository,
  );
  static final commits = RoutePath(
    path: '/commits',
    parent: global.RoutePaths.repository,
  );
  static final settings = RoutePath(
    path: '/settings',
    parent: global.RoutePaths.repository,
  );
}

class Routes {
  final List<RouteDefinition> all = [
    RouteDefinition(
      routePath: RoutePaths.files,
      component: FilesComponentNgFactory,
      useAsDefault: true,
    ),
    RouteDefinition(
      routePath: RoutePaths.commits,
      component: CommitsComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.settings,
      component: SettingsComponentNgFactory,
    ),
  ];
}
