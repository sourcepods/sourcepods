import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/login/login_component.template.dart';
import 'package:gitpods/src/user/user_list_component.template.dart';
import 'package:gitpods/src/user/user_profile_component.template.dart';

class RoutePaths {
  static final root = RoutePath(path: '/');
  static final login = RoutePath(path: '/login');
  static final userList = RoutePath(path: '/users');
  static final userProfile = RoutePath(path: '/:username');
  static final settings = RoutePath(path: '/settings/...');
  static final repository = RoutePath(path: '/:owner/:name/...');
  static final repositoryCreate = RoutePath(path: '/new');
  static final notFound = RoutePath(path: '/**');
}

class Routes {
  final List<RouteDefinition> all = [
    RouteDefinition.redirect(
      routePath: RoutePaths.root,
      redirectTo: RoutePaths.userList.toUrl(),
    ),
    RouteDefinition(
      routePath: RoutePaths.login,
      component: LoginComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.userList,
      component: UserListComponentNgFactory,
    ),
    RouteDefinition(
      routePath: RoutePaths.userProfile,
      component: UserProfileComponentNgFactory,
    )
  ];
}


//@RouteConfig(const [
//  const Redirect(
//    path: '/',
//    redirectTo: const ['UserList'],
//  ),
//  const Route(
//    path: '/users',
//    name: 'UserList',
//    component: UserListComponent,
//  ),
//  const Route(
//    path: '/login',
//    name: 'Login',
//    component: LoginComponent,
//  ),
//  const Route(
//    path: '/:username',
//    name: 'UserProfile',
//    component: UserProfileComponent,
//  ),
//  const Route(
//    path: '/settings/...',
//    name: 'Settings',
//    component: SettingsComponent,
//  ),
//  const Route(
//    path: '/:owner/:name/...',
//    name: 'Repository',
//    component: RepositoryComponent,
//  ),
//  const Route(
//    path: '/new',
//    name: 'RepositoryCreate',
//    component: RepositoryCreateComponent,
//  ),
//  const Route(
//    path: '/**',
//    name: 'NotFound',
//    component: NotFoundComponent,
//  )
//])
