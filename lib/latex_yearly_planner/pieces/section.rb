# frozen_string_literal: true

module LatexYearlyPlanner
  module Pieces
    class Section
      def initialize(config, section_config, header, body)
        @config = config
        @section_config = section_config
        @header = header
        @body = body
      end

      def enabled?
        section_config.enabled || false
      end

      def generate
        NotImplementedError
      end

      private

      attr_reader :config, :section_config, :header, :body
    end
  end
end
