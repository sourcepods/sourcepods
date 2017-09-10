import 'package:angular/angular.dart';
import 'package:gitpods/app_component.dart';
import 'package:http/browser_client.dart';

void main() {
  bootstrap(AppComponent, [
    provide(BrowserClient, useFactory: () => new BrowserClient(), deps: [])
  ]);
}
