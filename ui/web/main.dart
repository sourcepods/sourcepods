import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';
import 'package:sourcepods/api.dart';
import 'package:sourcepods/app_component.template.dart' as ng;

import 'main.template.dart' as self;

@GenerateInjector([
  routerProviders,
  ClassProvider(Client, useClass: BrowserClient),
  ClassProvider(Location, useClass: Location),
  ClassProvider(API),
])
final InjectorFactory injector = self.injector$Injector;

void main() {
  runApp(ng.AppComponentNgFactory, createInjector: injector);
}
