import 'package:angular_router/angular_router.dart';
import 'package:gitpods/routes.dart' as _parent;
import 'package:gitpods/src/repository/commits/commits_component.template.dart';
import 'package:gitpods/src/repository/files/files_component.template.dart';
import 'package:gitpods/src/repository/settings/settings_component.template.dart';

class RoutePaths {
  static final commits = RoutePath(
    path: 'commits',
    parent: _parent.RoutePaths.repository,
  );
  static final settings = RoutePath(
    path: 'settings',
    parent: _parent.RoutePaths.repository,
  );
  static final files = RoutePath(
    path: '',
    parent: _parent.RoutePaths.repository,
    useAsDefault: true,
  );
}

class Routes {
  static final List<RouteDefinition> all = [
    RouteDefinition(
      routePath: RoutePaths.commits,
      component: CommitsComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.settings,
      component: SettingsComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.files,
      component: FilesComponentNgFactory,
      useAsDefault: true,
    ),
  ];
}
