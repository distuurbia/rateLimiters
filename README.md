# About the project:
## It presents realized rate limiters algorythms, that is easy to use in the middleware.
## 1. Fixed window
### Example of implementation in the middleware:
   ```  
   func (cm *CustomMiddleware) MiddlewareFixedWindow(next echo.HandlerFunc) echo.HandlerFunc {
  	  return func(c echo.Context) error {
  		  res := cm.fw.CheckIfRequestAllowed()
  		  if !res {
  			  return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later fixedWindow")
  		  }
  		  return next(c)
  	  }
    }
   ```
## 2. Token bucket
### Example of implementation in the middleware:
   ```
   func (cm *CustomMiddleware) MiddlewareTokenBucket(next echo.HandlerFunc) echo.HandlerFunc {
  	  return func(c echo.Context) error {
  		  const tokens = 100
  		  res := cm.tb.CheckIfRequestAllowed(tokens)
  		  if !res {
  			  return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later tokenBucket")
  		  }
  		  return next(c)
  	  }
    }
   ```
## 3. Sliding log
### Example of implementation in the middleware:
   ```
   func (cm *CustomMiddleware) MiddlewareSlidingLog(next echo.HandlerFunc) echo.HandlerFunc {
      return func(c echo.Context) error {
         const (
            maxRequests        = 10
            intervalToBeParsed = "10s"
         )
         interval, err := time.ParseDuration(intervalToBeParsed)
         if err != nil {
            err = fmt.Errorf("CustomMiddleware-time.ParseDuration-err: %w", err)
            slog.Error(err.Error())
            return echo.NewHTTPError(http.StatusInternalServerError)
         }
         res := cm.sl.CheckIfRequestAllowed(c.RealIP(), uuid.NewString(), interval, maxRequests)
         if !res {
            return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later sliding log")
         }
         return next(c)
      }
   }
   ```
## 4. Sliding window
### Example of implementation in the middleware:
   ```
      func (cm *CustomMiddleware) MiddlewareSlidingWindow(next echo.HandlerFunc) echo.HandlerFunc {
      	return func(c echo.Context) error {
      		const (
      			maxRequests        = 10
      			intervalToBeParsed = "10s"
      		)
      		interval, err := time.ParseDuration(intervalToBeParsed)
      		if err != nil {
      			err = fmt.Errorf("CustomMiddleware-time.ParseDuration-err: %w", err)
      			slog.Error(err.Error())
      			return echo.NewHTTPError(http.StatusInternalServerError)
      		}
      		res, err := cm.sw.CheckIfRequestAllowed(c.RealIP(), interval, maxRequests)
      		if err != nil {
      			err = fmt.Errorf("CustomMiddleware-cm.sw.CheckIfRequestAllowed-err: %w", err)
      			slog.Error(err.Error())
      			return echo.NewHTTPError(http.StatusBadRequest)
      		}
      		if !res {
      			return echo.NewHTTPError(http.StatusTooManyRequests, "Try again later sliding window")
      		}
      		return next(c)
      	}
      }
   ```
