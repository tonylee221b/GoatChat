import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:goat_chat/src/app.dart';
import 'package:goat_chat/src/features/chat/presentation/chat_home_page.dart';

void main() {
  testWidgets('shows nickname gate before entering chat home', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    expect(find.text('Before you start'), findsOneWidget);
    expect(find.text('Enter your nickname'), findsOneWidget);
    expect(find.widgetWithText(ElevatedButton, 'Continue'), findsOneWidget);
    expect(find.byType(ChatHomePage), findsNothing);
  });

  testWidgets('keeps continue disabled until nickname is valid', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    final continueButton = tester.widget<ElevatedButton>(
      find.widgetWithText(ElevatedButton, 'Continue'),
    );
    expect(continueButton.onPressed, isNull);

    await tester.enterText(find.byType(TextField), 'Go');
    await tester.pump();

    final enabledButton = tester.widget<ElevatedButton>(
      find.widgetWithText(ElevatedButton, 'Continue'),
    );
    expect(enabledButton.onPressed, isNotNull);
  });

  testWidgets('submits nickname and navigates directly to chat home', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    await tester.enterText(find.byType(TextField), 'Goat');
    await tester.pump();
    await tester.tap(find.widgetWithText(ElevatedButton, 'Continue'));
    await tester.pumpAndSettle();

    expect(find.byType(ChatHomePage), findsOneWidget);
    expect(find.text('GoatChat'), findsOneWidget);
    expect(find.text('All Chats'), findsOneWidget);
  });

  testWidgets('uses the existing monochrome palette across the flow', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    final context = tester.element(find.byType(Scaffold).first);
    expect(Theme.of(context).colorScheme.primary, const Color(0xFF111111));
    expect(Theme.of(context).scaffoldBackgroundColor, const Color(0xFFF7F7F5));
  });
}
