import 'package:flutter/material.dart';

abstract final class AppColors {
  static const ink = Color(0xFF111111);
  static const textPrimary = Color(0xFF151320);
  static const border = Color(0xFF191526);
  static const warmWhite = Color(0xFFF7F7F5);
  static const mutedText = Color(0xFF8C8C86);
  static const tabTrack = Color(0xFFEDEDE9);
  static const tabText = Color(0xFF73736D);
}

abstract final class AppTheme {
  static ThemeData get light {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: AppColors.ink,
        primary: AppColors.ink,
        surface: AppColors.warmWhite,
      ),
      scaffoldBackgroundColor: AppColors.warmWhite,
      fontFamily: 'Roboto',
    );
  }
}
