import 'package:flutter/material.dart';
import 'package:frontend/screens/activities_screen.dart';
import 'package:frontend/screens/activity_screen.dart';
import 'screens/welcome_screen.dart';
import 'screens/login_screen.dart';
import 'screens/registration_screen.dart';
import 'screens/home_screen.dart';
import 'screens/nutrition_service_screens/add_food_screen.dart';
import 'screens/nutrition_service_screens/add_meal_screen.dart';
import 'screens/nutrition_service_screens/add_water_screen.dart';
import 'screens/nutrition_service_screens/foods_screen.dart';
import 'screens/nutrition_service_screens/meal_screen.dart';
import 'screens/nutrition_service_screens/water_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Healthy Summer',
      theme: ThemeData(primarySwatch: Colors.orange),
      initialRoute: '/',
      routes: {
        '/': (context) => const WelcomeScreen(),
        '/login': (context) => const LoginScreen(),
        '/register': (context) => const RegistrationScreen(),
        '/home': (context) => const HomeScreen(),
        '/activities': (context) => const ActivitiesScreen(),
        '/add_activity': (context) => const AddActivityScreen(),
        '/meals': (context) => const MealsScreen(),
        '/add_meal': (context) => const AddMealScreen(),
        '/water': (context) => const WaterScreen(),
        '/add_water': (context) => const AddWaterScreen(),
        '/foods': (context) => const FoodsScreen(),
        '/add_food': (context) => const AddFoodScreen(),
      },
    );
  }
}
