import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/src/login/login_component.template.dart';
import 'package:sourcepods/src/not_found_component.template.dart';
import 'package:sourcepods/src/repository/repository_component.template.dart';
import 'package:sourcepods/src/repository/repository_create_component.template.dart';
import 'package:sourcepods/src/settings/settings_component.template.dart';
import 'package:sourcepods/src/user/user_list_component.template.dart';
import 'package:sourcepods/src/user/user_profile_component.template.dart';

class RoutePaths {
  static final login = RoutePath(path: 'login');
  static final repositoryCreate = RoutePath(path: 'new');
  static final settings = RoutePath(path: 'settings');
  static final userList = RoutePath(path: 'users');
  static final userProfile = RoutePath(path: ':username');
  static final repository = RoutePath(path: ':owner/:name');
  static final root = RoutePath(path: '');
  static final notFound = RoutePath(path: '.+');
}

class Routes {
  static final List<RouteDefinition> all = [
    RouteDefinition(
      routePath: RoutePaths.login,
      component: LoginComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.repositoryCreate,
      component: RepositoryCreateComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.settings,
      component: SettingsComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.userList,
      component: UserListComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.userProfile,
      component: UserProfileComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.repository,
      component: RepositoryComponentNgFactory,
    ),
    RouteDefinition.redirect(
      routePath: RoutePaths.root,
      redirectTo: RoutePaths.userList.toUrl(),
    ),
    RouteDefinition(
      routePath: RoutePaths.notFound,
      component: NotFoundComponentNgFactory,
    ),
  ];
}
