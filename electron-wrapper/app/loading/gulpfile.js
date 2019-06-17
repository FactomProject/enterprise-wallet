var gulp = require('gulp');
var postcss = require('gulp-postcss');
var autoprefixer = require('autoprefixer');
var cssnext = require('postcss-cssnext');
var precss = require('precss');
var del = require('del');
var runSequence = require('run-sequence');
var livereload = require('gulp-livereload');

gulp.task('css', function () {
  var processors = [
    cssnext,
    precss,
  ];
  return gulp.src('./src/styles/index.css')
    .pipe(postcss(processors))
    .pipe(postcss([ autoprefixer() ]))
    .pipe(gulp.dest('css'))
    .pipe(livereload());
});

gulp.task('watch', ['css'], function() {
    livereload.listen();
    gulp.watch('src/**/*.css', ['css']);
});

gulp.task('default', ['watch']);

